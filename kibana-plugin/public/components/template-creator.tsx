import React, { useState } from 'react';
import uuid from 'uuid';

import {
  EuiButton,
  EuiButtonEmpty,
  EuiFieldText,
  EuiForm,
  EuiFormRow,
  EuiModal,
  EuiModalHeader,
  EuiModalHeaderTitle,
  EuiModalBody,
  EuiModalFooter,
  EuiOverlayMask,
  EuiSelect,
  EuiSelectOption,
  EuiTextArea,
} from '@elastic/eui';

import { Protocol } from '../../common/checks/protocol';

import { ITemplate, protocolFromString } from '../../common/types';

interface TemplateCreatorProps {
  onClose: () => void;
  onCreate: (template: ITemplate) => void;
}

function optionsFromProtocol(): EuiSelectOption[] {
  // Get a list of string values from the members of the Protocol enum
  const protocolValues: string[] = Object.values(Protocol).filter(x => typeof x === 'string');

  // Create the option objects from the list
  const protocolOptions: EuiSelectOption[] = [];
  protocolValues.forEach((proto) => {
    protocolOptions.push({
      value: proto,
      text: proto,
    });
  })

  return protocolOptions;
}

export function TemplateCreator(props: TemplateCreatorProps) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [protocol, setProtocol] = useState(Protocol.Noop);

  function createTemplate() {
    props.onCreate({
      id: uuid.v4(),
      title: title,
      description: description,
      protocol: protocol,
    });
  }

  return (
    <EuiOverlayMask>
      <EuiModal onClose={props.onClose}>
        <EuiModalHeader>
          <EuiModalHeaderTitle>Create a new template</EuiModalHeaderTitle>
        </EuiModalHeader>
        <EuiModalBody>
          <EuiForm>
            <EuiFormRow label="Title">
              <EuiFieldText name="title" value={title} onChange={event => setTitle(event.target.value)} />
            </EuiFormRow>
            <EuiFormRow label="Description">
              <EuiTextArea name="description" value={description} onChange={event => setDescription(event.target.value)} />
            </EuiFormRow>
            <EuiFormRow label="Protocol">
              <EuiSelect options={optionsFromProtocol()} value={protocol} onChange={(event) => {
                setProtocol(protocolFromString(event.target.value));
              }} />
            </EuiFormRow>
          </EuiForm>
        </EuiModalBody>
        <EuiModalFooter>
          <EuiButtonEmpty onClick={props.onClose}>Cancel</EuiButtonEmpty>
          <EuiButton onClick={createTemplate} fill>Save</EuiButton>
        </EuiModalFooter>
      </EuiModal>
    </EuiOverlayMask>
  )
}