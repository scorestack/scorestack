import React, { Fragment, useState } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';

import {
  EuiPage,
  EuiPageBody,
  EuiPageContent,
} from '@elastic/eui';

import { CoreStart } from '../../../../src/core/public';
import { NavigationPublicPluginStart } from '../../../../src/plugins/navigation/public';

import { PLUGIN_ID } from '../../common';
import { ITemplate } from '../../common/types';
import { Protocol } from '../../common/checks/protocol';

import { NoTemplatePrompt } from './no-template-prompt';
import { TemplateCreator } from './template-creator';
import { TemplateTable } from './template-table';

interface ScoreStackAppProps {
  basename: string;
  notifications: CoreStart['notifications'];
  http: CoreStart['http'];
  navigation: NavigationPublicPluginStart;
}

const startingTemplates: ITemplate[] = [
  {
    id: '0001',
    title: 'Wordpress - Twenty Twenty',
    description:
      'Checks the content of the index page for the Wordpress default Twenty Twenty theme.',
    protocol: Protocol.HTTP,
  },
];

export const ScoreStackApp = (props: ScoreStackAppProps) => {
  const [templates, setTemplates] = useState(startingTemplates);
  const [showingCreator, setShowingCreator] = useState(false);

  function onCreateTemplate() {
    setShowingCreator(true);
  }

  function onCloseTemplateCreator() {
    setShowingCreator(false);
  }

  function saveNewTemplate(template: ITemplate) {
    setTemplates(templates.concat(template));
    onCloseTemplateCreator();
  }

  function renderTable(items: ITemplate[]): React.ReactNode {
    // If there are no items, instead render an EuiEmptyPrompt
    if (items.length === 0) {
      return <NoTemplatePrompt onClick={onCreateTemplate} />;
    } else {
      return (
        <TemplateTable
          basename={props.basename}
          items={items}
          onCreateTemplate={onCreateTemplate}
          addToast={(toast) => {
            return props.notifications.toasts.add(toast);
          }}
        />
      );
    }
  }

  let creator: React.ReactNode;

  if (showingCreator) {
    creator = <TemplateCreator onClose={onCloseTemplateCreator} onCreate={saveNewTemplate} />
  }

  // Render the application DOM.
  return (
    <Router basename={props.basename}>
      <Fragment>
        <props.navigation.ui.TopNavMenu appName={PLUGIN_ID} />
        {/* TODO: make page resize to be smaller when displaying an empty prompt */}
        <EuiPage restrictWidth="1000px">
          <EuiPageBody>
            <EuiPageContent>{renderTable(templates)}</EuiPageContent>
          </EuiPageBody>
        </EuiPage>
        {creator}
      </Fragment>
    </Router>
  );
};
