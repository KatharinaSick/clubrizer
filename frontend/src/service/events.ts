import axios from '@/plugins/axios'

export interface Category {
  id: string
  name: string
  color: string
  picture: string
  sortOrder: number
  customLabel?: string
}

export interface Creator {
  id: string
  givenName: string
  familyName: string
  nickName: string
  picture?: string
}

export interface Event {
  id: string
  title: string
  description: string
  startTime: Date
  location: string
  categoryId: string
  creator: Creator
}

export interface EventAttendee {
  id: string
  givenName: string
  familyName: string
  nickName: string
  picture?: string
  response: boolean
}

export interface EventResponses {
  going: number
  notGoing: number
  currentUserResponse: boolean | null
  attendees: EventAttendee[]
}

export interface EventDetail extends Event {
  category: Category
  responses: EventResponses
  canDelete?: boolean
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

export const upsertEventResponse = async (eventId: string, response: boolean): Promise<void> => {
  await axios.put(`/events/${eventId}/response`, { response })
}

export const deleteEvent = async (eventId: string): Promise<void> => {
  await axios.delete(`/events/${eventId}`)
}
