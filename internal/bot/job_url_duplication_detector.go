// SPDX-License-Identifier: GPL-2.0
package bot

import (
	"encoding/json"
	"fmt"
	netUrl "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/justincampbell/timeago"
	redis "github.com/novitoll/novitoll_daemon_bot/pkg/redis_client"
	"github.com/sirupsen/logrus"
	"mvdan.cc/xurls"
)

func JobUrlDuplicationDetector(j *Job) (interface{}, error) {
	msg := j.req.Message
	f := j.app.Features.UrlDuplication

	if !f.Enabled {
		return false, nil
	}

	urls := xurls.Relaxed().FindAllString(msg.Text, -1)

	if len(urls) == 0 {
		return false, nil
	}

	redisConn := redis.GetRedisConnection()
	defer redisConn.Close()

	for i, url := range urls {
		// this should be controlled in JobAdDetector
		if strings.Contains(url, "t.me/") {
			continue
		}

		j.app.Logger.Info(fmt.Sprintf("Checking %d/%d URL - %s",
			i+1, len(urls), url))

		if f.IgnoreHostnames {
			u, err := netUrl.ParseRequestURI(url)
			if err != nil || u.Path == "" {
				continue
			}
		}

		redisKey := fmt.Sprintf("%s-%d-%s", REDIS_USER_SENT_URL, msg.Chat.Id, strings.ToLower(url))

		jsonStr, _ := redisConn.Get(redisKey).Result()

		if jsonStr != "" {
			j.app.Logger.WithFields(logrus.Fields{
				"chat": msg.Chat.Id,
				"url":  url,
			}).Warn("This message contains the duplicate URL")

			var duplicatedMsg BotInReqMsg

			json.Unmarshal([]byte(jsonStr), &duplicatedMsg)

			_, err := j.onURLDuplicate(&duplicatedMsg)
			if err != nil {
				return false, err
			}
		} else {
			payload, err := json.Marshal(msg)

			// should not be the case here
			if err != nil {
				j.app.Logger.Fatal("Can not marshal BotInReq.Message from Redis")
				return false, err
			}

			err2 := redisConn.Set(redisKey, payload,
				time.Duration(f.RelevanceTimeout)*time.Second).Err()

			if err2 != nil {
				j.app.Logger.WithFields(logrus.Fields{
					"err": err,
				}).Fatal("[-] Can not put the message to Redis")
				return false, err
			}
		}
	}

	// check if user is allowed to post URLs
	k := fmt.Sprintf("%s-%d-%d", REDIS_USER_VERIFIED, msg.Chat.Id, msg.From.Id)

	t0, _ := redisConn.Get(j).Result()
	if t0 == "" {
		j.app.Logger("User is not in verified map. Probably user is aight.")
		return false, nil
	}

	t0i, err := strconv.Atoi(t0.(string))
	if err != nil {
		j.app.Logger.Warn("Could not convert string to int")
		return false, err
	}

	if t0i > 0 && (time.Now().Unix()-int64(t0i) <= NEWCOMER_URL_POST_DELAY) {
		// new users can not post URLs until NEWCOMER_URL_POST_DELAY mins
		n := int(NEWCOMER_URL_POST_DELAY / 60)

		botReply := fmt.Sprintf(j.app.Features.NewcomerQuestionnare.
			I18n[j.app.Lang].AuthMessageURLPost, n)

		admins := j.app.ChatAdmins[msg.Chat.Id]

		botReply += fmt.Sprintf(" CC: @%s, %s", BDFL, strings.Join(admins, ", "))

		go func() {
			select {
			case <-time.After(time.Duration(TIME_TO_DELETE_REPLY_MSG) * time.Second):
				j.DeleteMessage(&msg)
			}
		}()

		return j.SendMessage(botReply, msg.MessageId)
	}

	return nil, nil
}

func (j *Job) onURLDuplicate(duplicatedMsg *BotInReqMsg) (interface{}, error) {
	msg := j.req.Message
	f := j.app.Features.UrlDuplication

	j.app.Logger.Info("POST HTTP request on duplicate detection")

	t := time.Since(time.Unix(duplicatedMsg.Date, 0))
	d, _ := time.ParseDuration(t.String())

	botReply := fmt.Sprintf(f.I18n[j.app.Lang].WarnMessage,
		duplicatedMsg.From.Username, timeago.FromDuration(d))

	reply, err := j.SendMessage(botReply, msg.MessageId)
	if err != nil {
		return false, err
	}

	if reply != nil {
		// cleanup reply messages
		go func() {
			select {
			case <-time.After(time.Duration(TIME_TO_DELETE_REPLY_MSG+10) * time.Second):
				go j.DeleteMessage(reply)
			}
		}()
	}

	return reply, err
}
