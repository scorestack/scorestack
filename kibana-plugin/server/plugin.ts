import {
  PluginInitializerContext,
  CoreSetup,
  CoreStart,
  Plugin,
  Logger,
} from '../../../src/core/server';

import { PLUGIN_ID, PLUGIN_NAME } from '../common';

import { ScoreStackPluginSetup, ScoreStackPluginStart } from './types';
import { defineRoutes } from './routes';

export class ScoreStackPlugin implements Plugin<ScoreStackPluginSetup, ScoreStackPluginStart> {
  private readonly logger: Logger;

  constructor(initializerContext: PluginInitializerContext) {
    this.logger = initializerContext.logger.get();
  }

  public setup(core: CoreSetup) {
    this.logger.debug(`${PLUGIN_ID}: Setup`);
    const router = core.http.createRouter();

    // Register server side APIs
    defineRoutes(router);

    return {};
  }

  public start(core: CoreStart) {
    this.logger.debug(`${PLUGIN_ID}: Started`);
    return {};
  }

  public stop() {
    this.logger.debug(`${PLUGIN_ID}: Stopped`);
  }
}
