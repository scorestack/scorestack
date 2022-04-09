import { AppMountParameters, CoreSetup, CoreStart, Plugin } from 'kibana/public';
import { ScorestackPluginSetup, ScorestackPluginStart, AppPluginStartDependencies } from './types';
import { PLUGIN_NAME } from '../common';

export class ScorestackPlugin implements Plugin<ScorestackPluginSetup, ScorestackPluginStart> {
  public setup(core: CoreSetup): ScorestackPluginSetup {
    // Register an application into the side navigation menu
    core.application.register({
      id: 'scorestack',
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

  public start(core: CoreStart): ScorestackPluginStart {
    return {};
  }

  public stop() {}
}
