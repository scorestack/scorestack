import { AppMountParameters, CoreSetup, CoreStart, Plugin } from '../../../src/core/public';

import { PLUGIN_ID, PLUGIN_NAME } from '../common';

import { ScoreStackPluginSetup, ScoreStackPluginStart, AppPluginStartDependencies } from './types';

export class ScoreStackPlugin implements Plugin<ScoreStackPluginSetup, ScoreStackPluginStart> {
  public setup(core: CoreSetup): ScoreStackPluginSetup {
    // Register an application into the side navigation menu
    core.application.register({
      id: `${PLUGIN_ID}`,
      title: PLUGIN_NAME,
      async mount(params: AppMountParameters) {
        // Load application bundle
        const { renderApp } = await import('./application');
        // Get start services as specified in kibana.json
        const [coreStart, depsStart] = await core.getStartServices();
        // Render the application
        return renderApp(coreStart, depsStart as AppPluginStartDependencies, params);
      },
    });

    // Return methods that should be available to other plugins
    return {};
  }

  public start(core: CoreStart): ScoreStackPluginStart {
    return {};
  }
}
