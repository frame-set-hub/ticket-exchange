import { useState, useEffect, useCallback } from 'react';
import type { Ticket, TicketFilters, CreateTicketParams } from '../../../domains/ticket/entities/Ticket';
import { getServices } from '../../../infrastructure/services/ServiceContainer';

export function useTickets() {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);
  const { ticketRepository } = getServices();

  const fetchTickets = useCallback(async (filters?: TicketFilters) => {
    setLoading(true);
    try {
      const data = await ticketRepository.getAll(filters);
      setTickets(data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchTickets();
  }, [fetchTickets]);

  return { tickets, loading, fetchTickets };
}

export function useMyTickets() {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);
  const { ticketRepository } = getServices();

  const fetchMyTickets = useCallback(async () => {
    setLoading(true);
    try {
      const data = await ticketRepository.getMyTickets();
      setTickets(data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  const deleteTicket = async (id: number) => {
    if (!confirm('Are you sure you want to delete this ticket?')) return;
    try {
      await ticketRepository.delete(id);
      fetchMyTickets();
    } catch (err: any) {
      alert(err.response?.data?.error || 'Failed to delete ticket');
    }
  };

  useEffect(() => {
    fetchMyTickets();
  }, [fetchMyTickets]);

  return { tickets, loading, deleteTicket };
}

export function useCreateTicket() {
  const [loading, setLoading] = useState(false);
  const { ticketRepository } = getServices();

  const createTicket = async (params: CreateTicketParams) => {
    setLoading(true);
    try {
      await ticketRepository.create(params);
    } finally {
      setLoading(false);
    }
  };

  return { createTicket, loading };
}
