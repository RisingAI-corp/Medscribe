// SampleComponent.js
import React from 'react';

interface Props {
  title: string;
  subtitle: string;
  backgroundColor: string;
  textColor: string;
  isDisabled: boolean;
}

const SampleComponent = ({ title, subtitle, backgroundColor, textColor, isDisabled }: Props) => {
  return (
    <div
      style={{
        backgroundColor,
        color: textColor,
        padding: '20px',
        border: '1px solid #ddd',
      }}
    >
      <h1>{title}</h1>
      <p>{subtitle}</p>
      <button disabled={isDisabled}>Click me now!</button>
    </div>
  );
};

export default SampleComponent;
