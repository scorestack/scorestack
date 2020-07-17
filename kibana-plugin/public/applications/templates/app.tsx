import { EuiPage, EuiPageBody } from '@elastic/eui';
import React, { useReducer } from 'react';
import { HashRouter, Route, Switch } from 'react-router-dom';
import uuid from 'uuid';
import { HttpStart, NotificationsStart } from '../../../../../src/core/public';
import { ITemplate } from '../../../common/types';
import { Home } from './home';
import { Template } from './template';
import { ActionType, TemplateAction } from './types';

interface TemplatesAppProps {
  basepath: string;
  http: HttpStart;
  notifs: NotificationsStart;
}

function templateReducer(state: ITemplate[], action: TemplateAction): ITemplate[] {
  switch (action.type) {
    case ActionType.Copy:
      return state.concat({
        id: uuid(),
        ...action.template,
      });
    case ActionType.Remove:
      return state.filter((template) => template.id !== action.template.id);
    case ActionType.Save:
      return state.map((template) => {
        return template.id === action.template.id ? action.template : template;
      });
  }
}

export function TemplatesApp(props: TemplatesAppProps): React.ReactElement {
  const [templates, dispatchTemplates] = useReducer(templateReducer, []);

  function copyTemplate(template: ITemplate): void {
    dispatchTemplates({ type: ActionType.Copy, template });
  }

  function removeTemplate(template: ITemplate): void {
    dispatchTemplates({ type: ActionType.Remove, template });
  }

  function saveTemplate(template: ITemplate): void {
    dispatchTemplates({ type: ActionType.Save, template });
  }

  // TODO: remove this
  function getTemplate(id: string): ITemplate {
    const result = templates.filter((template) => template.id === id);
    // TODO: do we really need to check the length here?
    return result.length === 1 ? result[0] : null;
  }

  return (
    <HashRouter basename={props.basepath}>
      <EuiPage restrictWidth="1000px">
        <EuiPageBody>
          <Switch>
            <Route exact path="/">
              <Home basepath={props.basepath} templates={templates} copyTemplate={copyTemplate} />
            </Route>
            <Route path="/template/:id">
              <Template
                get={getTemplate}
                copy={copyTemplate}
                remove={removeTemplate}
                save={saveTemplate}
              />
            </Route>
          </Switch>
        </EuiPageBody>
      </EuiPage>
    </HashRouter>
  );
}
