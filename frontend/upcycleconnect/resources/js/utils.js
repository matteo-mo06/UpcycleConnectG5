const ROLE_LABELS = {
    admin:   'Administrateur',
    artisan: 'Artisan',
    salarie: 'Salarié',
    user:    'Particulier',
}

const CONDITION_LABELS = {
    neuf:     'Neuf',
    bon_etat: 'Bon état',
    use:      'Usé',
}

export function formatDate(d) {
    if (!d) return '—'
    return new Date(d).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' })
}

export function conditionLabel(v) {
    return CONDITION_LABELS[v] ?? v
}

export function roleLabel(name) {
    return ROLE_LABELS[name] ?? (name ? name.charAt(0).toUpperCase() + name.slice(1) : 'Particulier')
}

export function primaryRole(roles = []) {
    for (const r of ['admin', 'artisan', 'salarie', 'user']) {
        if (roles.includes(r)) return ROLE_LABELS[r]
    }
    return 'Particulier'
}

export function fullName(user) {
    return `${user?.first_name ?? ''} ${user?.last_name ?? ''}`.trim() || 'Inconnu'
}

export function initials(user) {
    const f = user?.first_name?.[0] ?? ''
    const l = user?.last_name?.[0] ?? ''
    return (f + l).toUpperCase() || 'U'
}
