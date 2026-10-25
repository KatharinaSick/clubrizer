export const combineDateAndTime = (date: string, time: string): Date => {
  return new Date(`${date}T${time}`)
}

// toDateInput and toTimeInput split a date into the local "YYYY-MM-DD" and "HH:MM" values
// that date and time inputs expect, the reverse of combineDateAndTime.
export const toDateInput = (date: Date): string => {
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`
}

export const toTimeInput = (date: Date): string => {
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${pad(date.getHours())}:${pad(date.getMinutes())}`
}
