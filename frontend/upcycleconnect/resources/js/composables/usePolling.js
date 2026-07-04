import { onMounted, onUnmounted } from "vue";

export function usePolling(fn, interval = 2000, condition = () => true) {
    let timer = null;

    function start() {
        stop();
        timer = setInterval(() => {
            if (condition()) fn(true);
        }, interval);
    }

    function stop() {
        if (timer) {
            clearInterval(timer);
            timer = null;
        }
    }

    onMounted(() => {
        fn(false);
        start();
    });
    onUnmounted(stop);

    return { start, stop };
}
