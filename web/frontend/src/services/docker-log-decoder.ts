// Docker 的非 TTY 输出包含 8 字节流头，WebSocket 消息边界不一定与流帧对齐。
export class DockerLogDecoder {
  private pending = new Uint8Array(0)
  private remaining = 0
  private decoder = new TextDecoder()
  constructor(private tty: boolean) {}

  push(input: ArrayBuffer): string {
    if (this.tty) return this.decoder.decode(input, { stream: true })
    const bytes = new Uint8Array(input)
    const combined = new Uint8Array(this.pending.length + bytes.length)
    combined.set(this.pending)
    combined.set(bytes, this.pending.length)
    this.pending = combined
    let offset = 0
    let output = ''
    while (offset < this.pending.length) {
      if (this.remaining === 0) {
        if (this.pending.length - offset < 8) break
        const header = this.pending.subarray(offset, offset + 8)
        if (header[0] > 2 || header[1] !== 0 || header[2] !== 0 || header[3] !== 0) {
          output += this.decoder.decode(this.pending.subarray(offset), { stream: true })
          offset = this.pending.length
          break
        }
        this.remaining = new DataView(header.buffer, header.byteOffset, 8).getUint32(4)
        offset += 8
        if (!this.remaining) continue
      }
      const count = Math.min(this.remaining, this.pending.length - offset)
      output += this.decoder.decode(this.pending.subarray(offset, offset + count), { stream: true })
      offset += count
      this.remaining -= count
    }
    this.pending = this.pending.slice(offset)
    return output
  }
}

