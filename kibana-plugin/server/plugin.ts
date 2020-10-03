import {
  PluginInitializerContext,
  CoreSetup,
  CoreStart,
  Plugin,
  Logger,
} from '../../../src/core/server';

import { ScorestackPluginSetup, ScorestackPluginStart } from './types';
import { defineRoutes } from './routes';

export class ScorestackPlugin implements Plugin<ScorestackPluginSetup, ScorestackPluginStart> {
  private readonly logger: Logger;

  constructor(initializerContext: PluginInitializerContext) {
    this.logger = initializerContext.logger.get();
  }

  public setup(core: CoreSetup) {
    this.logger.debug('scorestack: Setup');
    const router = core.http.createRouter();

    // Register server side APIs
    defineRoutes(router, core.elasticsearch.legacy.client);

    return {};
  }

  public start(core: CoreStart) {
    this.logger.debug('scorestack: Started');
    return {};
  }

  public stop() { }
}
