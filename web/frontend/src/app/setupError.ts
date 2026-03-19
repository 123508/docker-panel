import { ElMessage } from "element-plus"

export function setupGlobalErrorHandler() {
    window.addEventListener("unhandledrejection", (event) => {
        const err = event.reason
        event.preventDefault()

        handleError(err)
    })

    window.addEventListener("error", (event) => {
        handleError(event.error)
    })
}

function handleError(err: any) {
    if (err?.handled) return

    const msg = err?.message || "系统错误"

    ElMessage.error(msg)

    err.handled = true
}