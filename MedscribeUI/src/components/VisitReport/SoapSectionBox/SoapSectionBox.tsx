import React from 'react';
import { Textarea } from '@mantine/core';
import { Accordion } from '@mantine/core';
import { CopyButton, Button } from '@mantine/core';
import { IconDownload } from '@tabler/icons-react';

interface SoapEditableSectionProps {
  title: string;
  text: string;
  onTextChange: (newText: string) => void;
}

const SoapSectionBox = React.memo(({ title, text, onTextChange }: SoapEditableSectionProps) => {
  return (
    <>
      <Accordion>
        <Accordion.Item key={title} value={title}>
          <Accordion.Control>
            <span style={{ fontSize: '1.25rem', fontWeight: 'bold' }}>{title}</span>
          </Accordion.Control>
          <Accordion.Panel>
            <Textarea
              value={text}
              placeholder="Input placeholder"
              onChange={event => onTextChange(event.currentTarget.value)} // Update text dynamically
            />
            <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '1rem' }}>
              <CopyButton value={text}>
                {({ copied, copy }) => (
                  <Button color={copied ? 'teal' : 'blue'} onClick={copy} rightSection={<IconDownload size={20} />}>
                    {copied ? 'Copied' : 'Copy'}
                  </Button>
                )}
              </CopyButton>
            </div>
          </Accordion.Panel>
        </Accordion.Item>
      </Accordion>
    </>
  );
});

export default SoapSectionBox;
