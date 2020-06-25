import React, { Fragment, useState } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';

import {
  EuiBasicTable,
  EuiBasicTableColumn,
  EuiButton,
  EuiEmptyPrompt,
  EuiLink,
  EuiPage,
  EuiPageBody,
  EuiPageContent,
  EuiPageContentBody,
  EuiPageContentHeader,
  EuiPageContentHeaderSection,
  EuiTitle,
} from '@elastic/eui';

import { CoreStart } from '../../../../src/core/public';
import { NavigationPublicPluginStart } from '../../../../src/plugins/navigation/public';

import { PLUGIN_ID, PLUGIN_NAME } from '../../common';
import { Protocol } from '../../common/checks/protocol';

interface ScoreStackAppProps {
  basename: string;
  notifications: CoreStart['notifications'];
  http: CoreStart['http'];
  navigation: NavigationPublicPluginStart;
}

interface Template {
  id: string;
  title: string;
  description: string;
  protocol: Protocol;
}

const startingTemplates: Template[] = [
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

  function createVisualizationClickHandler() {
    props.notifications.toasts.addInfo('Added template');
  }

  function editVisualizationClickHandler(item: Template) {
    props.notifications.toasts.addInfo(`Editing template: ${item.title}`);
  }

  function copyVisualizationClickHandler(item: Template) {
    props.notifications.toasts.addInfo(`Copied template: ${item.title}`);
  }

  function renderTitle(item: Template): React.ReactNode {
    return <EuiLink href={`${props.basename}/${item.id}`}>{item.title}</EuiLink>;
  }

  function renderTable(
    items: Template[],
    columns: Array<EuiBasicTableColumn<Template>>
  ): React.ReactNode {
    // If there are no items, instead render an EuiEmptyPrompt
    if (items.length === 0) {
      return (
        <EuiEmptyPrompt
          iconType="document"
          title={<h2>Create your first check template</h2>}
          body={
            <Fragment>
              <p>
                A template is used to configure most of the parameters for a check protocol to make
                it easy to create new checks.
              </p>
              <p>
                Templates can have attributes, which are parameters that are configured at runtime.
                Attributes allow you to reuse the same template for multiple checks if you only need
                to change a few simple things, like an IP address or a username.
              </p>
              <p>
                Once you create a template, you can start adding checks for the template by
                configuring values for the template&rsquo;s attributes.
              </p>
            </Fragment>
          }
          actions={
            <EuiButton fill onClick={createVisualizationClickHandler} iconType="plusInCircle">
              Create new template
            </EuiButton>
          }
        />
      );
    } else {
      return (
        <Fragment>
          <EuiPageContentHeader>
            <EuiPageContentHeaderSection>
              <EuiTitle>
                <h1>Check Templates</h1>
              </EuiTitle>
            </EuiPageContentHeaderSection>
            <EuiPageContentHeaderSection>
              <EuiButton fill onClick={createVisualizationClickHandler} iconType="plusInCircle">
                Create template
              </EuiButton>
            </EuiPageContentHeaderSection>
          </EuiPageContentHeader>
          <EuiPageContentBody>
            <EuiBasicTable
              items={items}
              columns={columns}
              tableLayout="auto"
              noItemsMessage="No templates found."
            />
          </EuiPageContentBody>
        </Fragment>
      );
    }
  }

  const columns: Array<EuiBasicTableColumn<Template>> = [
    {
      name: 'Title',
      render: renderTitle,
    },
    {
      field: 'protocol',
      name: 'Protocol',
    },
    {
      field: 'description',
      name: 'Description',
    },
    {
      name: 'Actions',
      actions: [
        {
          name: 'Edit',
          description: 'Edit Template',
          onClick: editVisualizationClickHandler,
          type: 'icon',
          icon: 'pencil',
        },
        {
          name: 'Copy',
          description: 'Copy Template',
          onClick: copyVisualizationClickHandler,
          type: 'icon',
          icon: 'copy',
        },
      ],
    },
  ];

  // Render the application DOM.
  return (
    <Router basename={props.basename}>
      <Fragment>
        <props.navigation.ui.TopNavMenu appName={PLUGIN_ID} />
        {/* TODO: make page resize to be smaller when displaying an empty prompt */}
        <EuiPage restrictWidth="1000px">
          <EuiPageBody>
            <EuiPageContent>{renderTable(templates, columns)}</EuiPageContent>
          </EuiPageBody>
        </EuiPage>
      </Fragment>
    </Router>
  );
};
