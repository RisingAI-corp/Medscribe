import landingBackground from '../../../assets/landing-bg.png';
import { Accordion, Container, Title } from '@mantine/core';
import { landingContent } from '../landingContent';

function LandingFAQ() {
  const faqData = landingContent.faq;

  return (
    <div className="h-full w-full flex relative justify-center items-center">
      
      <Container size="lg" className="py-10 w-full max-w-5xl mx-auto z-10">
        <Title order={2} className="mb-6 text-center text-gray-800">Frequently Asked Questions</Title>
        <Accordion
          variant="filled"
          radius="md"
          classNames={{
            item: 'border-b border-gray-200 last:border-b-0 bg-gray-50 mb-2 rounded-md shadow-sm',
            control: 'p-4 hover:bg-gray-100 transition-colors',
            label: 'font-medium text-gray-700'
          }}
        >
          {faqData.map((faq, index) => (
            <Accordion.Item key={index} value={`item-${index}`}>
              <Accordion.Control>{faq.question}</Accordion.Control>
              <Accordion.Panel>{faq.answer}</Accordion.Panel>
            </Accordion.Item>
          ))}
        </Accordion>
      </Container>
    </div>
  );
}

export default LandingFAQ;