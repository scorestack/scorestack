import React, { useState } from 'react';

import {
  EuiButton,
  EuiButtonEmpty,
  EuiModal,
  EuiModalHeader,
  EuiModalHeaderTitle,
  EuiModalBody,
  EuiModalFooter,
  EuiOverlayMask,
} from '@elastic/eui';

import { Protocol } from '../../common/checks/protocol';

import { ITemplate } from '../../common/types';

interface TemplateCreatorProps {
  onClose: () => void;
  onCreate: (template: ITemplate) => void;
}

export function TemplateCreator(props: TemplateCreatorProps) {
  function createTemplate() {
    props.onCreate({ id: '1234', title: 'Null', description: 'go away', protocol: Protocol.Noop });
  }

  return (
    <EuiOverlayMask>
      <EuiModal onClose={props.onClose}>
        <EuiModalHeader>
          <EuiModalHeaderTitle>Create a new template</EuiModalHeaderTitle>
        </EuiModalHeader>
        <EuiModalBody>
          <p>Some filler text...</p>
        </EuiModalBody>
        <EuiModalFooter>
          <EuiButtonEmpty onClick={props.onClose}>Cancel</EuiButtonEmpty>
          <EuiButton onClick={createTemplate} fill>Save</EuiButton>
        </EuiModalFooter>
      </EuiModal>
    </EuiOverlayMask>
  )
}