import type { ITicketRepository } from '../../domains/ticket/repositories/TicketRepository';
import type { Ticket, TicketFilters, CreateTicketParams } from '../../domains/ticket/entities/Ticket';
import apiClient from './apiClient';

export class TicketRepository implements ITicketRepository {
  async getAll(filters?: TicketFilters): Promise<Ticket[]> {
    const params = new URLSearchParams();
    if (filters?.title) params.append('title', filters.title);
    if (filters?.category) params.append('category', filters.category);
    if (filters?.min_price) params.append('min_price', filters.min_price);
    if (filters?.max_price) params.append('max_price', filters.max_price);

    const { data } = await apiClient.get(`/tickets?${params.toString()}`);
    return data;
  }

  async getMyTickets(): Promise<Ticket[]> {
    const { data } = await apiClient.get('/tickets/my');
    return data;
  }

  async create(params: CreateTicketParams): Promise<Ticket> {
    const { data } = await apiClient.post('/tickets', params);
    return data;
  }

  async delete(id: number): Promise<void> {
    await apiClient.delete(`/tickets/${id}`);
  }
}
