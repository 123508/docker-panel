import { ElMessageBox } from 'element-plus'

// 敏感操作统一在请求发送前确认，取消或关闭时不执行操作。
export async function confirmAction(action: string, target: string, consequence: string): Promise<boolean> {
  try {
    await ElMessageBox.confirm(`确定${action} ${target}？${consequence}`, `${action}确认`, {
      confirmButtonText: `确认${action}`,
      cancelButtonText: '取消',
      type: 'warning',
      distinguishCancelAndClose: true
    })
    return true
  } catch {
    return false
  }
}
