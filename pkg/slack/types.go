package slack

// Response is the all messages returned by the query
type Response struct {
	Ok       bool      `json:"ok"`
	Messages []Message `json:"messages"`
}

// Payload contains slack API parameters
type Payload struct {
	token     string
	channel   string
	count     int     //Number of messages to return, between 1 and 1000.
	inclusive bool    //Include messages with latest or oldest timestamp in results.
	latest    float64 //Optional, default=now, End of time range of messages to include in results.
	oldest    string  //Optional, default=0 Start of time range of messages to include in results.
	unreads   bool    //Optional, default=0 Include unread_count_display in the output?
	text      string
}

// Message is the individual message response
type Message struct {
	ClientMsgID string `json:"client_msg_id"`
	Type        string `json:"type"`
	Text        string `json:"text"`
	User        string `json:"user"`
	Ts          int    `json:"ts"`
	Team        string `json:"team"`
}

// PostSlackMessage is the struct for posting messages to Slack
type PostSlackMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// UserObject is the struct for user data
type UserObject struct {
	ID       string `json:"id"`
	TeamID   string `json:"team_id"`
	Name     string `json:"name"`
	Deleted  bool   `json:"deleted"`
	Color    string `json:"color"`
	RealName string `json:"real_name"`
	Tz       string `json:"tz"`
	TzLabel  string `json:"tz_label"`
	TzOffset int    `json:"tz_offset"`
	Profile  struct {
		AvatarHash            string `json:"avatar_hash"`
		StatusText            string `json:"status_text"`
		StatusEmoji           string `json:"status_emoji"`
		RealName              string `json:"real_name"`
		DisplayName           string `json:"display_name"`
		RealNameNormalized    string `json:"real_name_normalized"`
		DisplayNameNormalized string `json:"display_name_normalized"`
		Email                 string `json:"email"`
		ImageOriginal         string `json:"image_original"`
		Image24               string `json:"image_24"`
		Image32               string `json:"image_32"`
		Image48               string `json:"image_48"`
		Image72               string `json:"image_72"`
		Image192              string `json:"image_192"`
		Image512              string `json:"image_512"`
		Team                  string `json:"team"`
	} `json:"profile"`
	IsAdmin           bool `json:"is_admin"`
	IsOwner           bool `json:"is_owner"`
	IsPrimaryOwner    bool `json:"is_primary_owner"`
	IsRestricted      bool `json:"is_restricted"`
	IsUltraRestricted bool `json:"is_ultra_restricted"`
	IsBot             bool `json:"is_bot"`
	Updated           int  `json:"updated"`
	IsAppUser         bool `json:"is_app_user"`
	Has2Fa            bool `json:"has_2fa"`
}
