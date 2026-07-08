#!/usr/bin/env bash
set -euo pipefail

# Reinitialise la base a partir d'un fichier de dump (scenario) et relance le backend.
# A executer depuis la racine du repo, sur le serveur web.
#
# Usage :
#   ./db/apply_scenario.sh dump_vide.txt
#   ./db/apply_scenario.sh dump_plein.txt

DUMP_FILE="${1:?Usage: $0 <chemin_vers_dump>}"
ENV_FILE="upcycle_connect-api/.env"
COMPOSE_FILE="docker-compose.prod.yml"

if [ ! -f "$DUMP_FILE" ]; then
  echo "Fichier introuvable : $DUMP_FILE" >&2
  exit 1
fi

if [ ! -f "$ENV_FILE" ]; then
  echo "$ENV_FILE introuvable - lance ce script depuis la racine du repo." >&2
  exit 1
fi

DB_HOST=$(grep -E '^DB_HOST=' "$ENV_FILE" | cut -d= -f2-)
DB_PORT=$(grep -E '^DB_PORT=' "$ENV_FILE" | cut -d= -f2-)
DB_NAME=$(grep -E '^DB_NAME=' "$ENV_FILE" | cut -d= -f2-)
DB_USER=$(grep -E '^DB_USER=' "$ENV_FILE" | cut -d= -f2-)

echo "Cible : $DB_USER@$DB_HOST:$DB_PORT/$DB_NAME"
read -r -s -p "Mot de passe MySQL pour $DB_USER : " DB_PASSWORD
echo
read -r -p "Ceci va SUPPRIMER puis recreer '$DB_NAME' avec le contenu de '$DUMP_FILE'. Continuer ? [y/N] " CONFIRM
if [[ "$CONFIRM" != "y" && "$CONFIRM" != "Y" ]]; then
  echo "Annule."
  exit 1
fi

MYSQL="mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD"

echo "[1/3] Recreation de la base '$DB_NAME'..."
$MYSQL -e "
DROP DATABASE IF EXISTS \`$DB_NAME\`;
CREATE DATABASE \`$DB_NAME\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
"

echo "[2/3] Import de '$DUMP_FILE'..."
mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" --default-character-set=utf8mb4 "$DB_NAME" < "$DUMP_FILE"

echo "[3/3] Redemarrage du backend..."
docker compose -f "$COMPOSE_FILE" restart backend

echo
echo "Termine. Suis les logs avec :"
echo "  docker compose -f $COMPOSE_FILE logs -f backend"
