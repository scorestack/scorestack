import React, { useState } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';

import {
  EuiButton,
  EuiHorizontalRule,
  EuiPage,
  EuiPageBody,
  EuiPageContent,
  EuiPageContentBody,
  EuiPageContentHeader,
  EuiPageHeader,
  EuiTitle,
  EuiText,
} from '@elastic/eui';

import { CoreStart } from '../../../../src/core/public';
import { NavigationPublicPluginStart } from '../../../../src/plugins/navigation/public';

import { PLUGIN_ID, PLUGIN_NAME } from '../../common';

interface ScoreStackAppProps {
  basename: string;
  notifications: CoreStart['notifications'];
  http: CoreStart['http'];
  navigation: NavigationPublicPluginStart;
}

export const ScoreStackApp = (props: ScoreStackAppProps) => {
  const [timestamp, setTimestamp] = useState<string | undefined>();

  const onClickHandler = () => {
    // Use the core http service to make a response to the server API.
    props.http.get('/api/scorestack/example').then((res) => {
      setTimestamp(res.time);
      // Use the core notifications service to display a success message.
      props.notifications.toasts.addSuccess('Data updated');
    });
  };

  // Render the application DOM.
  return (
    <Router basename={props.basename}>
      <>
        <props.navigation.ui.TopNavMenu appName={PLUGIN_ID} />
        <EuiPage restrictWidth="1000px">
          <EuiPageBody>
            <EuiPageHeader>
              <EuiTitle size="l">
                <h1>{PLUGIN_NAME}</h1>
              </EuiTitle>
            </EuiPageHeader>
            <EuiPageContent>
              <EuiPageContentHeader>
                <EuiTitle>
                  <h2>Congratulations, you have successfully created a new Kibana Plugin!</h2>
                </EuiTitle>
              </EuiPageContentHeader>
              <EuiPageContentBody>
                <EuiText>
                  <p>
                    Look through the generated code and check out the plugin development
                    documentation.
                  </p>
                  <EuiHorizontalRule />
                  <p>Last timestamp: {timestamp ? timestamp : 'Unknown'}</p>
                  <EuiButton type="primary" size="s" onClick={onClickHandler}>
                    Get data
                  </EuiButton>
                </EuiText>
              </EuiPageContentBody>
            </EuiPageContent>
          </EuiPageBody>
        </EuiPage>
      </>
    </Router>
  );
};
