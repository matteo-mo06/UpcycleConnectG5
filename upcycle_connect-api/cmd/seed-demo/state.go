package main

type payementPool struct {
	all             []string
	announcementPay []string
	subscriptionPay []string
	formationPay    []string
	eventPay        []string
	projectPay      []string
	buyerOf         map[string]string
	formationOf     map[string]string
}

type state struct {
	userIDs      []string
	artisanIDs   []string
	salarieIDs   []string
	regularIDs   []string
	userRoleName map[string]string

	notificationIDs []string

	adviceIDs []string

	projectIDs     []string
	projectCreator map[string]string
	projectDeleted map[string]bool

	formationIDs     []string
	formationCreator map[string]string
	formationDeleted map[string]bool

	eventIDs     []string
	eventCreator map[string]string

	announcementIDs       []string
	announcementOwner     map[string]string
	announcementCategory  map[string]string
	announcementState     map[string]string
	announcementBuyer     map[string]string
	announcementDeleted   map[string]bool
	announcementLockerIDs []string

	professionalRequestIDs []string

	topicIDs     []string
	topicAuthor  map[string]string
	topicDeleted map[string]bool

	postIDs     []string
	postTopic   map[string]string
	postAuthor  map[string]string
	postDeleted map[string]bool

	advertisementIDs []string

	subscriptionIDs  []string
	subscriptionUser map[string]string

	payements payementPool

	reportIDs       []string
	reportResolved  []string
	reportedUserFor map[string]string

	userSanctionIDs []string
}

func newState() *state {
	return &state{
		payements: payementPool{
			buyerOf:     make(map[string]string),
			formationOf: make(map[string]string),
		},
		userRoleName:         make(map[string]string),
		projectCreator:       make(map[string]string),
		projectDeleted:       make(map[string]bool),
		formationCreator:     make(map[string]string),
		formationDeleted:     make(map[string]bool),
		eventCreator:         make(map[string]string),
		announcementOwner:    make(map[string]string),
		announcementCategory: make(map[string]string),
		announcementState:    make(map[string]string),
		announcementBuyer:    make(map[string]string),
		announcementDeleted:  make(map[string]bool),
		topicAuthor:          make(map[string]string),
		topicDeleted:         make(map[string]bool),
		postTopic:            make(map[string]string),
		postAuthor:           make(map[string]string),
		postDeleted:          make(map[string]bool),
		subscriptionUser:     make(map[string]string),
		reportedUserFor:      make(map[string]string),
	}
}
