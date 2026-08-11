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
  cancelledAt?: string
}

export interface EventAttendee {
  id: string
  givenName: string
  familyName: string
  nickName: string
  picture?: string
  response: boolean
  isKid: boolean
  // parent is the owning parent's display name, present only for kid attendees.
  parent?: string
}

// MyResponse is one person the current account can RSVP for (the account holder, if
// they participate, plus each approved kid) with that person's current response.
export interface MyResponse {
  id: string
  isSelf: boolean
  name: string
  picture?: string | null
  response: boolean | null
}

export interface EventResponses {
  going: number
  notGoing: number
  attendees: EventAttendee[]
  myResponses: MyResponse[]
}

export interface EventDetail extends Event {
  category: Category
  responses: EventResponses
  canDelete?: boolean
  canCancel?: boolean
}

export interface CreateEventRequest {
  title: string
  description: string | null
  startTime: Date
  location: string
  categoryId: string
}

export const listPastEvents = async (): Promise<Event[]> => {
  const response = await axios.get('/events/past')
  return response.data ?? []
}

export const createEvent = async (event: CreateEventRequest): Promise<Event> => {
  const response = await axios.post('/events', event)
  return response.data
}

// upsertEventResponse sets a response for one person. Omit kidId for the account
// holder's own response, or pass a kid's id to respond for that kid.
export const upsertEventResponse = async (
  eventId: string,
  response: boolean,
  kidId?: string,
): Promise<void> => {
  await axios.put(`/events/${eventId}/response`, { response, kidId: kidId ?? null })
}

// Comment is a message an account holder left on an event, with the author's display
// details included so it can be rendered without a second lookup.
export interface Comment {
  id: string
  body: string
  createdAt: string
  author: Creator
}

export const listComments = async (eventId: string): Promise<Comment[]> => {
  const response = await axios.get(`/events/${eventId}/comments`)
  return response.data
}

export const createComment = async (eventId: string, body: string): Promise<Comment> => {
  const response = await axios.post(`/events/${eventId}/comments`, { body })
  return response.data
}

export const deleteEvent = async (eventId: string): Promise<void> => {
  await axios.delete(`/events/${eventId}`)
}

export const cancelEvent = async (eventId: string): Promise<void> => {
  await axios.post(`/events/${eventId}/cancel`)
}

export const uncancelEvent = async (eventId: string): Promise<void> => {
  await axios.post(`/events/${eventId}/uncancel`)
}
