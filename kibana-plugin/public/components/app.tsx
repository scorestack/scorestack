import {
  EuiPage,
  EuiPageBody,
  EuiPageContent,
  EuiPageContentBody,
  EuiPageContentHeader,
  EuiPageContentHeaderSection,
  EuiPageSideBar,
  EuiSideNav,
  EuiText,
  EuiTitle,
} from '@elastic/eui';
import React, { Fragment, useEffect, useState } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import { CoreStart } from '../../../../src/core/public';
import { CheckAttributes } from '../../common';
import { Attribute } from './attribute';

interface ScorestackAppProps {
  basename: string;
  http: CoreStart['http'];
}

const initialNavItems = [
  {
    name: 'Loading...',
    id: 'loading',
  },
];

export function ScorestackApp(props: ScorestackAppProps): React.ReactElement {
  const [navItems, setNavItems] = useState(initialNavItems);
  const [focusedContent, setFocusedContent] = useState<React.ReactElement>(
    <EuiText>Click on any of your checks to the left to configure their attributes.</EuiText>
  );

  useEffect(() => {
    props.http
      .get('/api/scorestack/attribute')
      .then((response: CheckAttributes) => {
        // Construct nav headers from the group names
        const newNavItems = Object.entries(response).map(([groupName, groupChecks]) => {
          // Construct nav items for each check within the group
          const subNavItems = Object.entries(groupChecks).map(([checkId, check]) => {
            return {
              name: check.name,
              id: checkId,
              onClick: () => {
                // Create the attributes that will be displayed on the page
                const attributes = Object.entries(check.attributes).map(([name, value]) => (
                  <Attribute
                    key={`${checkId}-${name}`}
                    id={checkId}
                    name={name}
                    value={value}
                    http={props.http}
                  />
                ));

                // Set the new center content
                setFocusedContent(
                  <Fragment>
                    <EuiPageContentHeader>
                      <EuiPageContentHeaderSection>
                        <EuiTitle>
                          <h2>{check.name}</h2>
                        </EuiTitle>
                      </EuiPageContentHeaderSection>
                    </EuiPageContentHeader>
                    <EuiPageContentBody>{attributes}</EuiPageContentBody>
                  </Fragment>
                );
              },
            };
          });

          return {
            name: groupName,
            id: groupName,
            items: subNavItems,
          };
        });

        setNavItems(newNavItems);
      })
      .catch((error) => {
        // TODO: handle this with a toast
        console.log('Promise rejected - failed to load attributes');
        console.log(error);
      });
  }, [props.http]);

  return (
    <Router basename={props.basename}>
      <EuiPage restrictWidth="1000px">
        <EuiPageSideBar>
          <EuiSideNav mobileTitle="Checks" items={navItems} />
        </EuiPageSideBar>
        <EuiPageBody>
          <EuiPageContent>{focusedContent}</EuiPageContent>
        </EuiPageBody>
      </EuiPage>
    </Router>
  );
}
