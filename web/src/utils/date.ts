export function nowDatetime() {
  return formatTime(new Date(), 'YYYY-MM-DD HH:mm:ss')
}

export function formatTime(date: Date, format: string): string {
  const map: { [key: string]: string } = {
    YYYY: date.getFullYear().toString(),
    MM: (date.getMonth() + 1).toString().padStart(2, '0'),
    DD: date.getDate().toString().padStart(2, '0'),
    HH: date.getHours().toString().padStart(2, '0'),
    mm: date.getMinutes().toString().padStart(2, '0'),
    ss: date.getSeconds().toString().padStart(2, '0'),
  }
  return format.replace(/YYYY|MM|DD|HH|mm|ss/g, (match) => map[match])
}

export function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  return dateStr.slice(0, 10)
}

export function formatDateTime(dateStr: string): string {
  if (!dateStr) return '-'
  return dateStr.slice(0, 19)
}

export function RangeDate(range: number): string {
  if (!range) return ''
  const d = new Date()
  d.setDate(d.getDate() + range)
  return formatTime(d, 'YYYY-MM-DD')
}
