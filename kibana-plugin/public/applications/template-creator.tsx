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

import { Protocol, ProtocolList } from '../../common/checks';

import { ITemplate, protocolFromString } from '../../common/types';

interface TemplateCreatorProps {
  onClose: () => void;
  onCreate: (template: ITemplate) => void;
}

function optionsFromProtocol(): EuiSelectOption[] {
  // Create the option objects from the list
  const protocolOptions: EuiSelectOption[] = [];
  ProtocolList.forEach((proto) => {
    protocolOptions.push({
      value: proto,
      text: proto,
    });
  });

  return protocolOptions;
}

export function TemplateCreator(props: TemplateCreatorProps) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [protocol, setProtocol] = useState(Protocol.Noop);

  function createTemplate() {
    props.onCreate({
      id: uuid.v4(),
      title,
      description,
      protocol,
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
              <EuiFieldText
                name="title"
                value={title}
                onChange={(event) => setTitle(event.target.value)}
              />
            </EuiFormRow>
            <EuiFormRow label="Description">
              <EuiTextArea
                name="description"
                value={description}
                onChange={(event) => setDescription(event.target.value)}
              />
            </EuiFormRow>
            <EuiFormRow label="Protocol">
              <EuiSelect
                options={optionsFromProtocol()}
                value={protocol}
                onChange={(event) => {
                  setProtocol(protocolFromString(event.target.value));
                }}
              />
            </EuiFormRow>
          </EuiForm>
        </EuiModalBody>
        <EuiModalFooter>
          <EuiButtonEmpty onClick={props.onClose}>Cancel</EuiButtonEmpty>
          <EuiButton onClick={createTemplate} fill>
            Save
          </EuiButton>
        </EuiModalFooter>
      </EuiModal>
    </EuiOverlayMask>
  );
}
