import { PluginSetupContract as FeaturesPluginSetup } from '../../../x-pack/plugins/features/server';

/* eslint-disable @typescript-eslint/no-empty-interface, prettier/prettier */
export interface ScoreStackPluginSetup { }
export interface ScoreStackPluginStart { }
/* eslint-enable @typescript-eslint/no-empty-interface, prettier/prettier */

export interface ScoreStackPluginDeps {
  features: FeaturesPluginSetup;
}
