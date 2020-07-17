import React from 'react';
import ReactDOM from 'react-dom';

import { AppMountParameters, HttpSetup, NotificationsStart } from '../../../../../src/core/public';

import { TemplatesApp } from './app';

export function render(params: AppMountParameters, http: HttpSetup, notifs: NotificationsStart) {
  ReactDOM.render(
    <TemplatesApp basepath={params.appBasePath} http={http} notifs={notifs} />,
    params.element
  );

  return () => ReactDOM.unmountComponentAtNode(params.element);
}
