import type { IAuthRepository } from '../../domains/auth/repositories/AuthRepository';
import type { ITicketRepository } from '../../domains/ticket/repositories/TicketRepository';
import type { ITransactionRepository } from '../../domains/transaction/repositories/TransactionRepository';
import type { IChatRepository } from '../../domains/chat/repositories/ChatRepository';
import { AuthRepository } from '../api/authRepository';
import { TicketRepository } from '../api/ticketRepository';
import { TransactionRepository } from '../api/transactionRepository';
import { ChatRepository } from '../api/chatRepository';

interface Services {
  authRepository: IAuthRepository;
  ticketRepository: ITicketRepository;
  transactionRepository: ITransactionRepository;
  chatRepository: IChatRepository;
}

let services: Services | null = null;

export function getServices(): Services {
  if (!services) {
    services = {
      authRepository: new AuthRepository(),
      ticketRepository: new TicketRepository(),
      transactionRepository: new TransactionRepository(),
      chatRepository: new ChatRepository(),
    };
  }
  return services;
}
