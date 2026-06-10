export function StringLike(keyStr: string, containerStr: string) {
  if (keyStr.length == 0) return true
  return containerStr.indexOf(keyStr) >= 0
}

export function isChinese(str: string): boolean {
  return /[\u4e00-\u9fa5]/.test(str)
}

export function isEnglish(str: string): boolean {
  return /^[A-Za-z\s]*$/.test(str)
}

export function tansErr(err: string, vName: string, iName: string): string {
  return err.replace(vName, iName)
}
