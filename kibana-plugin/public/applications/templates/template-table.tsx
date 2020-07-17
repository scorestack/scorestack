import {
  EuiBasicTableColumn,
  EuiButton,
  EuiLink,
  EuiPageContent,
  EuiPageContentBody,
  EuiPageContentHeader,
  EuiPageContentHeaderSection,
  EuiTitle,
} from '@elastic/eui';
import React from 'react';
import { EuiButtonIcon } from '@elastic/eui';
import { EuiBasicTable } from '@elastic/eui';
import { ITemplate } from '../../../common/types';

interface TemplateTableProps {
  basepath: string;
  items: ITemplate[];
  copyTemplate: (item: ITemplate) => void;
}

export function TemplateTable(props: TemplateTableProps): React.ReactElement {
  /* The render prop for EuiBasicTable custom actions isn't React.ReactNode;
  it's some other weird thing, so we're just gonna let TS infer the return type
  for this function. It doesn't really matter anyway */
  function renderEditButton(item: ITemplate) {
    return <EuiButtonIcon href={`${props.basepath}#/template/${item.id}`} iconType="pencil" />;
  }

  function renderTitle(item: ITemplate): React.ReactNode {
    return <EuiLink href={`${props.basepath}#/template/${item.id}`}>{item.title}</EuiLink>;
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
          render: renderEditButton,
        },
        {
          name: 'Copy',
          description: 'Copy Template',
          onClick: props.copyTemplate,
          type: 'icon',
          icon: 'copy',
        },
      ],
    },
  ];

  return (
    <EuiPageContent>
      <EuiPageContentHeader>
        <EuiPageContentHeaderSection>
          <EuiTitle>
            <h1>Check Templates</h1>
          </EuiTitle>
        </EuiPageContentHeaderSection>
        <EuiPageContentHeaderSection>
          <EuiButton fill href={`${props.basepath}#/template`} iconType="plusInCircle">
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
    </EuiPageContent>
  );
}
