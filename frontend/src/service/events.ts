import axios from '@/plugins/axios'

export interface Event {
  id: string
  title: string
  description: string
  startTime: Date
  location: string
  categoryId: string
}

export interface CreateEventRequest {
  title: string
  description: string | null
  startTime: Date
  location: string
  categoryId: string
}

export const createEvent = async (event: CreateEventRequest): Promise<Event> => {
  const response = await axios.post('/events', event)
  return response.data
}
