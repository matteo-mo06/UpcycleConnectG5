package main

import (
	"database/sql"
	"fmt"
)

type advertisementPlan struct {
	ID         int
	Weeks      int
	PriceCents int
}

type refs struct {
	categoryIDByName   map[string]string
	roleIDByName       map[string]string
	adminUserID        string
	lockerIDs          []string
	advertisementPlans []advertisementPlan
	subscriptionPlanID string
	commissionRate     float64
}

func loadRefs(db *sql.DB) (*refs, error) {
	r := &refs{
		categoryIDByName: make(map[string]string),
		roleIDByName:     make(map[string]string),
	}

	catRows, err := db.Query("SELECT id_category, name_category FROM category")
	if err != nil {
		return nil, err
	}
	defer catRows.Close()
	for catRows.Next() {
		var id, name string
		if err := catRows.Scan(&id, &name); err != nil {
			return nil, err
		}
		r.categoryIDByName[name] = id
	}
	for _, th := range themes {
		if _, ok := r.categoryIDByName[th.Category]; !ok {
			return nil, fmt.Errorf("catégorie introuvable en base : %s (les migrations ont-elles été appliquées ?)", th.Category)
		}
	}

	roleRows, err := db.Query("SELECT id_role, name_role FROM role")
	if err != nil {
		return nil, err
	}
	defer roleRows.Close()
	for roleRows.Next() {
		var id, name string
		if err := roleRows.Scan(&id, &name); err != nil {
			return nil, err
		}
		r.roleIDByName[name] = id
	}
	for _, needed := range []string{"admin", "artisan", "salarie", "user"} {
		if _, ok := r.roleIDByName[needed]; !ok {
			return nil, fmt.Errorf("rôle introuvable en base : %s", needed)
		}
	}

	err = db.QueryRow(`
		SELECT u.id_user FROM user u
		JOIN user_role ur ON ur.id_user = u.id_user
		WHERE ur.id_role = ?
		LIMIT 1`, r.roleIDByName["admin"],
	).Scan(&r.adminUserID)
	if err != nil {
		return nil, fmt.Errorf("aucun admin trouvé en base, exécutez d'abord cmd/seed : %w", err)
	}

	lockerRows, err := db.Query("SELECT id_locker FROM locker ORDER BY locker_number")
	if err != nil {
		return nil, err
	}
	defer lockerRows.Close()
	for lockerRows.Next() {
		var id string
		if err := lockerRows.Scan(&id); err != nil {
			return nil, err
		}
		r.lockerIDs = append(r.lockerIDs, id)
	}
	if len(r.lockerIDs) == 0 {
		return nil, fmt.Errorf("aucun locker trouvé en base")
	}

	planRows, err := db.Query("SELECT id, weeks, price_cents FROM advertisement_plans")
	if err != nil {
		return nil, err
	}
	defer planRows.Close()
	for planRows.Next() {
		var p advertisementPlan
		if err := planRows.Scan(&p.ID, &p.Weeks, &p.PriceCents); err != nil {
			return nil, err
		}
		r.advertisementPlans = append(r.advertisementPlans, p)
	}
	if len(r.advertisementPlans) == 0 {
		return nil, fmt.Errorf("aucun plan de publicité trouvé en base")
	}

	err = db.QueryRow("SELECT id_plan FROM subscription_plans WHERE is_active = 1 LIMIT 1").Scan(&r.subscriptionPlanID)
	if err != nil {
		return nil, fmt.Errorf("aucun plan d'abonnement actif trouvé en base : %w", err)
	}

	var rateStr string
	if err := db.QueryRow("SELECT value_setting FROM platform_settings WHERE key_setting = 'commission_rate'").Scan(&rateStr); err != nil {
		return nil, fmt.Errorf("commission_rate introuvable dans platform_settings : %w", err)
	}
	if _, err := fmt.Sscanf(rateStr, "%f", &r.commissionRate); err != nil {
		return nil, fmt.Errorf("commission_rate invalide (%q) : %w", rateStr, err)
	}

	return r, nil
}
