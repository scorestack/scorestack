import React, { Fragment, useState } from 'react';
import { HashRouter, Route, Switch, useParams } from 'react-router-dom';

import { EuiPage, EuiPageBody, EuiPageContent } from '@elastic/eui';

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

interface TemplateProps {
  templates: ITemplate[];
}

function Template({ templates }: TemplateProps) {
  const { id } = useParams();
  const tmpl = templates.filter((t) => {
    return t.id === id;
  })[0];
  return (
    <Fragment>
      <h1>{tmpl.title}</h1>
      <p>{tmpl.description}</p>
    </Fragment>
  );
}

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
    creator = <TemplateCreator onClose={onCloseTemplateCreator} onCreate={saveNewTemplate} />;
  }

  // Render the application DOM.
  return (
    <HashRouter basename={props.basename}>
      <EuiPage restrictWidth="1000px">
        <Switch>
          <Route exact path="/">
            {/* TODO: make page resize to be smaller when displaying an empty prompt */}
            <EuiPageBody>
              <EuiPageContent>{renderTable(templates)}</EuiPageContent>
            </EuiPageBody>
            {creator}
          </Route>
          <Route path="/:id">
            <Template templates={templates} />
          </Route>
        </Switch>
      </EuiPage>
    </HashRouter>
  );
};
