// SPDX-License-Identifier: GPL-2.0
package config

type FeaturesConfig struct {
	UrlDuplication       UrlConfigOptions               `json:"urlDuplication"`
	NewcomerQuestionnare NewcomerConfigOptions          `json:"newcomerQuestionnare"`
	StickersDetection    StickersDetectionConfigOptions `json:"stickersDetection"`
	MessageStatistics    MessageStatisticsConfigOptions `json:"messageStatistics"`
	Administration       AdministrationConfigOptions    `json:"administration"`
	AdDetection          AdDetectionConfigOptions       `json:"adDetection"`
}

type NewcomerConfigOptions struct {
	Enabled        bool  `json:"enabled"`
	ActionNotify   bool  `json:"actionNotify"`
	AuthTimeout    uint8 `json:"authTimeout"`
	KickBanTimeout uint8 `json:"kickBanTimeout"`
	I18n           map[string]struct {
		WelcomeMessage string `json:"welcomeMessage"`
		AuthOKMessage  string `json:"authOKMessage"`
		AuthMessage    string `json:"authMessage"`
		AuthMessageCached string `json:"authMessageCached"`
	} `json:"i18n"`
}

type UrlConfigOptions struct {
	Enabled          bool `json:"enabled"`
	IgnoreHostnames  bool `json:"ignoreHostnames"`
	ActionNotify     bool `json:"actionNotify"`
	RelevanceTimeout int  `json:"relevanceTimeout"`
	I18n             map[string]struct {
		WarnMessage string `json:"warnMessage"`
	} `json:"i18n"`
}

type StickersDetectionConfigOptions struct {
	Enabled bool `json:"enabled"`
	I18n    map[string]struct {
		WarnMessage string `json:"warnMessage"`
	}
}

type MessageStatisticsConfigOptions struct {
	Enabled bool `json:"enabled"`
	I18n    map[string]struct {
		WarnMessageTooFreq string `json:"warnMessageTooFreq"`
		WarnMessageTooLong string `json:"warnMessageTooLong"`
	} `json:"i18n"`
}

type AdministrationConfigOptions struct {
	LogLevel string   `json:"logLevel"`
	I18n     map[string]struct {
		CronJobNewcomersReport string   `json:"cronJobNewcomersReport"`
		CronJobUserMsgReport   string   `json:"cronJobUserMsgReport"`
	} `json:"i18n"`
}

type AdDetectionConfigOptions struct {
	Enabled bool `json:"enabled"`
	I18n    map[string]struct {
		WarnMessage string `json:"warnMessage"`
	} `json:"i18n"`
}
