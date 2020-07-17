import React from 'react';
import { ITemplate } from '../../../common/types';
import { HomeEmpty } from './home-empty';
import { TemplateTable } from './template-table';

interface HomeProps {
  basepath: string;
  templates: ITemplate[];
  copyTemplate: (item: ITemplate) => void;
}

export function Home(props: HomeProps): React.ReactElement {
  if (props.templates.length === 0) {
    return <HomeEmpty basepath={props.basepath} />;
  } else {
    return (
      <TemplateTable
        basepath={props.basepath}
        items={props.templates}
        copyTemplate={props.copyTemplate}
      />
    );
  }
}
