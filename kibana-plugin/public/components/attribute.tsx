import React, { useState } from 'react';
import {
  EuiFlexGroup,
  EuiFlexItem,
  EuiButtonIcon,
  EuiButton,
  EuiFormRow,
  EuiFieldText,
  EuiPopover,
  EuiText,
} from '@elastic/eui';
import { CoreStart } from '../../../../src/core/public';

interface AttributeProps {
  name: string;
  id: string;
  value: string;
  http: CoreStart['http'];
}

export function Attribute(props: AttributeProps): React.ReactElement {
  const [value, setValue] = useState(props.value);
  const [formValue, setFormValue] = useState('');
  const [visible, setVisible] = useState(false);
  const [loading, setLoading] = useState(false);

  const showButton = <EuiButtonIcon iconType="eye" onClick={() => setVisible(!visible)} />;

  function onSaveButtonClick() {
    // Start the loading spinner
    setLoading(true);

    props.http
      .post(`/api/scorestack/attribute/${props.id}/${props.name}`, {
        body: JSON.stringify({ value: formValue }),
        headers: {
          'Content-Type': 'application/json',
        },
      })
      .then(() => {
        setLoading(false);
        setValue(formValue);
        setFormValue('');
      })
      .catch((error) => {
        // TODO: handle this with a toast
        /* eslint-disable no-console */
        console.log('Promise rejected - failed to save attribute value');
        console.log(error);
        /* eslint-enable no-console */
      });
  }

  return (
    <EuiFlexGroup style={{ maxWidth: 600 }}>
      <EuiFlexItem>
        <EuiFormRow label={props.name}>
          <EuiFieldText value={formValue} onChange={(event) => setFormValue(event.target.value)} />
        </EuiFormRow>
      </EuiFlexItem>
      <EuiFlexItem grow={false}>
        <EuiFormRow hasEmptyLabelSpace>
          <EuiPopover
            id="showValue"
            ownFocus
            button={showButton}
            isOpen={visible}
            closePopover={() => setVisible(false)}
          >
            <EuiText>
              <code>{value}</code>
            </EuiText>
          </EuiPopover>
        </EuiFormRow>
      </EuiFlexItem>
      <EuiFlexItem grow={false}>
        <EuiFormRow hasEmptyLabelSpace>
          <EuiButton isLoading={loading} onClick={onSaveButtonClick}>
            Save
          </EuiButton>
        </EuiFormRow>
      </EuiFlexItem>
    </EuiFlexGroup>
  );
}
