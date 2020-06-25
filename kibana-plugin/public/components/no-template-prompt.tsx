import React, { Fragment } from 'react';

import { EuiButton, EuiEmptyPrompt } from '@elastic/eui';

interface NoTemplatePromptProps {
  onClick: () => void;
}

export function NoTemplatePrompt(props: NoTemplatePromptProps) {
  return (
    <EuiEmptyPrompt
      iconType="document"
      title={<h2>Create your first check template</h2>}
      body={
        <Fragment>
          <p>
            A template is used to configure most of the parameters for a check protocol to make it
            easy to create new checks.
          </p>
          <p>
            Templates can have attributes, which are parameters that are configured at runtime.
            Attributes allow you to reuse the same template for multiple checks if you only need to
            change a few simple things, like an IP address or a username.
          </p>
          <p>
            Once you create a template, you can start adding checks for the template by configuring
            values for the template&rsquo;s attributes.
          </p>
        </Fragment>
      }
      actions={
        <EuiButton fill onClick={props.onClick} iconType="plusInCircle">
          Create new template
        </EuiButton>
      }
    />
  );
}
