package main

type theme struct {
	Category      string
	Announcements []string
	Project       string
	Formation     string
	Advice        string
}

var themes = []theme{
	{
		Category:      "Couture & Textile",
		Announcements: []string{"Chutes de tissu", "Bobines de fil", "Boutons dépareillés"},
		Project:       "Coussins en tissu recyclé",
		Formation:     "Initiation couture",
		Advice:        "Bien choisir son fil",
	},
	{
		Category:      "Menuiserie & Bois",
		Announcements: []string{"Palette en bois", "Chutes de contreplaqué", "Planches de récup"},
		Project:       "Table basse en palette",
		Formation:     "Fabrication de mobilier en palettes",
		Advice:        "Reconnaître un bois traité",
	},
	{
		Category:      "Électronique",
		Announcements: []string{"Vieil ordinateur", "Câbles et chargeurs", "Pièces détachées"},
		Project:       "Réparation d'appareils électroniques",
		Formation:     "Bases de la réparation électronique",
		Advice:        "Démonter un appareil en sécurité",
	},
	{
		Category:      "Maroquinerie",
		Announcements: []string{"Chutes de cuir", "Vieux sac à restaurer"},
		Project:       "Sac en cuir recyclé",
		Formation:     "Initiation à la maroquinerie",
		Advice:        "Nettoyer un cuir récupéré",
	},
	{
		Category:      "Céramique & Poterie",
		Announcements: []string{"Vaisselle ébréchée", "Argile d'occasion"},
		Project:       "Mosaïque en céramique",
		Formation:     "Restauration de céramique",
		Advice:        "Réparer une pièce fissurée",
	},
	{
		Category:      "Jardinage & Compostage",
		Announcements: []string{"Bidons plastique", "Palettes pour compost"},
		Project:       "Composteur en palettes",
		Formation:     "Initiation au compostage",
		Advice:        "Équilibrer son compost",
	},
	{
		Category:      "Peinture & Décoration",
		Announcements: []string{"Chaises à repeindre", "Restes de peinture"},
		Project:       "Relooking d'un meuble",
		Formation:     "Peinture sur mobilier",
		Advice:        "Préparer une surface avant peinture",
	},
	{
		Category:      "Récupération Métaux",
		Announcements: []string{"Barres de fer", "Touret de câble"},
		Project:       "Table en touret de câble",
		Formation:     "Atelier soudure débutant",
		Advice:        "Manipuler le métal en sécurité",
	},
	{
		Category:      "Papeterie & Carton",
		Announcements: []string{"Cartons d'emballage", "Chutes de papier"},
		Project:       "Étagère en carton",
		Formation:     "Créations en carton",
		Advice:        "Renforcer une structure carton",
	},
	{
		Category:      "Mode & Couture Upcyclée",
		Announcements: []string{"Vieux jeans", "Vêtements usagés"},
		Project:       "Sac en jean recyclé",
		Formation:     "Upcycling textile",
		Advice:        "Choisir la bonne pièce à détourner",
	},
}

var announcementDescriptionSuffixes = []string{
	"bon état, à venir chercher",
	"à récupérer rapidement",
	"idéal pour bricoleur",
	"comme neuf",
	"lot à saisir",
	"disponible de suite",
	"quantité limitée",
	"état correct, petits défauts",
}

var conditionValues = []string{"Neuf", "Très bon état", "Bon état", "Usagé"}

var eventTitles = []string{
	"Atelier réparation",
	"Vide-dressing solidaire",
	"Marché des créateurs",
	"Conférence économie circulaire",
}

var topicTitles = []string{
	"Astuce nettoyage bois",
	"Retrouver des palettes près de Montreuil",
	"Retour d'expérience meuble",
	"Question réparation vélo",
}

var postBodies = []string{
	"Merci pour l'astuce",
	"J'ai eu le même souci",
	"Ça marche bien chez moi",
	"Bravo pour le résultat",
	"Quelqu'un a un contact ?",
}

var reportReasons = []string{
	"Contenu trompeur sur l'état du produit",
	"Prix incorrect indiqué",
	"Doublon d'annonce existante",
	"Description inexacte",
	"Propos inappropriés",
	"Contenu hors sujet de la plateforme",
}

var reportActionsTaken = []string{
	"Avertissement envoyé",
	"Suppression du contenu",
	"Aucune action nécessaire",
	"Suspension du compte",
}

var sanctionReasons = []string{
	"Non-respect des règles de la communauté",
	"Signalements répétés d'autres utilisateurs",
	"Annonce frauduleuse confirmée",
	"Comportement inapproprié sur le forum",
}
