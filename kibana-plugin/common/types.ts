import { Protocol } from './checks/protocol';

export interface ITemplate {
  id: string;
  title: string;
  description: string;
  protocol: Protocol;
}
