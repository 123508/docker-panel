import { reactive } from 'vue'

export interface ActionDialogState {
  visible: boolean
  title: string
  content: string
  okText: string
  isRunning: boolean
}

export interface RunWithDialogOptions {
  title: string
  pendingText: string
  successText?: string | (() => string)
  failureText?: string | ((err: any) => string)
}

export interface RunResult<T> {
  ok: boolean
  result?: T
  error?: any
}

export function useActionDialog() {
  const dialog = reactive<ActionDialogState>({
    visible: false,
    title: '',
    content: '',
    okText: '关闭',
    isRunning: false
  })

  // 打开弹窗，执行 fn，并把结果（成功/失败）反馈到弹窗上
  const runWithDialog = async <T>(opts: RunWithDialogOptions, fn: () => Promise<T>): Promise<RunResult<T>> => {
    dialog.title = opts.title
    dialog.content = opts.pendingText
    dialog.okText = '关闭'
    dialog.isRunning = true
    dialog.visible = true

    try {
      const result = await fn()
      const text =
        typeof opts.successText === 'function'
          ? opts.successText()
          : opts.successText || '✅ 操作成功'
      dialog.content = text
      return { ok: true, result }
    } catch (e: any) {
      const text =
        typeof opts.failureText === 'function'
          ? opts.failureText(e)
          : opts.failureText || `❌ 操作失败: \n${e?.message || '未知错误'}`
      dialog.content = text
      return { ok: false, error: e }
    } finally {
      dialog.isRunning = false
    }
  }

  return { dialog, runWithDialog }
}
