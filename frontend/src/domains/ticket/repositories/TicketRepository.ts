import { Ticket, TicketFilters, CreateTicketParams } from '../entities/Ticket';

export interface ITicketRepository {
  getAll(filters?: TicketFilters): Promise<Ticket[]>;
  getMyTickets(): Promise<Ticket[]>;
  create(params: CreateTicketParams): Promise<Ticket>;
  delete(id: number): Promise<void>;
}
