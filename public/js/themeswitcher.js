(() => {
    try {
        const stored = localStorage.getItem("themeMode");
        if (
            stored
                ? stored === "dark"
                : matchMedia("(prefers-color-scheme: dark)").matches
        ) {
            document.documentElement.classList.add("dark");
        }
    } catch (_) {}

    const apply = (dark) => {
        document.documentElement.classList.toggle("dark", dark);
        try {
            localStorage.setItem("themeMode", dark ? "dark" : "light");
        } catch (_) {}
    };

    document.addEventListener("basecoat:theme", (event) => {
        const mode = event.detail?.mode;
        apply(
            mode === "dark"
                ? true
                : mode === "light"
                  ? false
                  : !document.documentElement.classList.contains("dark"),
        );
    });
})();
