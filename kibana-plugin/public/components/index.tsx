import React from 'react';
import ReactDOM from 'react-dom';

import { AppMountParameters, HttpSetup, NotificationsStart } from '../../../../src/core/public';

import { TemplatesApp } from './templates';

export function render(params: AppMountParameters, http: HttpSetup, notifs: NotificationsStart) {
  ReactDOM.render(
    <TemplatesApp basename={params.appBasePath} http={http} notifs={notifs} />,
    params.element
  );

  return () => ReactDOM.unmountComponentAtNode(params.element);
}
