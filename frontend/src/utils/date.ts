export const combineDateAndTime = (date: string, time: string): Date => {
  return new Date(`${date}T${time}`)
}
