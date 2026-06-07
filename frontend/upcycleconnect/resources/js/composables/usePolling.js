import { onMounted, onUnmounted } from 'vue'

/**
 * Lance un polling automatique sur une fonction async.
 * Le premier appel se fait immédiatement, les suivants sont silencieux (pas de loading).
 *
 * @param {(silent: boolean) => Promise<void>} fn  - fonction à appeler (reçoit silent=true lors des rafraîchissements)
 * @param {number} interval  - intervalle en ms (défaut 2000)
 * @param {() => boolean} condition  - condition optionnelle pour suspendre le polling (ex: onglet inactif)
 */
export function usePolling(fn, interval = 2000, condition = () => true) {
    let timer = null

    function start() {
        stop()
        timer = setInterval(() => {
            if (condition()) fn(true)
        }, interval)
    }

    function stop() {
        if (timer) { clearInterval(timer); timer = null }
    }

    onMounted(() => { fn(false); start() })
    onUnmounted(stop)

    return { start, stop }
}
