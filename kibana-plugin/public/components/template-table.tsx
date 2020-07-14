import React, { Fragment } from 'react';

import {
  EuiBasicTable,
  EuiBasicTableColumn,
  EuiButton,
  EuiLink,
  EuiPageContentBody,
  EuiPageContentHeader,
  EuiPageContentHeaderSection,
  EuiTitle,
} from '@elastic/eui';

import { Toast, ToastInput } from '../../../../src/core/public/notifications';

import { ITemplate } from '../../common/types';

interface TemplateTableProps {
  basename: string;
  items: ITemplate[];
  onCreateTemplate: () => void;
  addToast: (toast: ToastInput) => Toast;
}

export function TemplateTable(props: TemplateTableProps) {
  function renderTitle(item: ITemplate): React.ReactNode {
    return <EuiLink href={`${props.basename}#/${item.id}`}>{item.title}</EuiLink>;
  }

  function onEditTemplate(item: ITemplate): void {
    props.addToast(`Editing template: ${item.title}`);
  }

  function onCopyTemplate(item: ITemplate): void {
    props.addToast(`Copied template: ${item.title}`);
  }

  const columns: Array<EuiBasicTableColumn<ITemplate>> = [
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
          onClick: onEditTemplate,
          type: 'icon',
          icon: 'pencil',
        },
        {
          name: 'Copy',
          description: 'Copy Template',
          onClick: onCopyTemplate,
          type: 'icon',
          icon: 'copy',
        },
      ],
    },
  ];

  return (
    <Fragment>
      <EuiPageContentHeader>
        <EuiPageContentHeaderSection>
          <EuiTitle>
            <h1>Check Templates</h1>
          </EuiTitle>
        </EuiPageContentHeaderSection>
        <EuiPageContentHeaderSection>
          <EuiButton fill onClick={props.onCreateTemplate} iconType="plusInCircle">
            Create template
          </EuiButton>
        </EuiPageContentHeaderSection>
      </EuiPageContentHeader>
      <EuiPageContentBody>
        <EuiBasicTable
          items={props.items}
          columns={columns}
          tableLayout="auto"
          noItemsMessage="No templates found."
        />
      </EuiPageContentBody>
    </Fragment>
  );
}
