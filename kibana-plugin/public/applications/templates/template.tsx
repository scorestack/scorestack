import {
  EuiButtonEmpty,
  EuiFieldText,
  EuiFlexGroup,
  EuiFlexItem,
  EuiForm,
  EuiFormRow,
  EuiPageContent,
  EuiPageContentBody,
  EuiPageContentHeader,
  EuiPageHeader,
  EuiSelect,
  EuiTextArea,
  EuiTitle,
} from '@elastic/eui';
import React, { Fragment, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Protocol, ProtocolList } from '../../../common/checks';
import { ITemplate } from '../../../common/types';

interface TemplateProps {
  get: (id: string) => ITemplate; // TODO: remove this
  copy: (template: ITemplate) => void;
  remove: (template: ITemplate) => void;
  save: (template: ITemplate) => void;
}

const initialTemplate: ITemplate = {
  id: '',
  title: '',
  description: '',
  protocol: Protocol.Noop,
};

export function Template(props: TemplateProps): React.ReactElement {
  // Get the right template object
  const { id } = useParams();
  const [template, setTemplate] = useState(props.get(id) ?? initialTemplate);
  // TODO: use the saved objects API instead of retrieving the template from parent
  // TODO: if a template ID doesn't exist, return to the templates homepage

  return (
    <Fragment>
      <EuiPageHeader>
        <EuiFlexGroup>
          <EuiFlexItem>
            <EuiButtonEmpty onClick={() => props.save(template)}>Save</EuiButtonEmpty>
          </EuiFlexItem>

          <EuiFlexItem>
            <EuiButtonEmpty onClick={() => props.remove(template)}>Remove</EuiButtonEmpty>
          </EuiFlexItem>

          <EuiFlexItem>
            <EuiButtonEmpty onClick={() => props.save(template)}>Copy</EuiButtonEmpty>
          </EuiFlexItem>
        </EuiFlexGroup>
      </EuiPageHeader>
      <EuiPageContent>
        <EuiPageContentHeader>
          <EuiTitle>
            <h1>Template Editor</h1>
          </EuiTitle>
        </EuiPageContentHeader>
        <EuiPageContentBody>
          <EuiForm>
            <EuiFormRow label="Title">
              <EuiFieldText
                name="title"
                value={template.title}
                onChange={(event) => setTemplate({ title: event.target.value, ...template })}
              />
            </EuiFormRow>

            <EuiFormRow label="Description">
              <EuiTextArea
                name="description"
                value={template.description}
                onChange={(event) => setTemplate({ description: event.target.value, ...template })}
              />
            </EuiFormRow>

            <EuiFormRow label="Protocol">
              <EuiSelect
                options={ProtocolList.map((proto) => {
                  return { value: proto, text: proto };
                })}
                value={template.protocol}
                onChange={(event) => setTemplate({ protocol: event.target.value, ...template })}
              />
            </EuiFormRow>
          </EuiForm>
        </EuiPageContentBody>
      </EuiPageContent>
    </Fragment>
  );
}
