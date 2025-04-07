import { SearchResultItem } from '../FollowUpSearchModal/SearchResults/SearchResults';

const PatientBackgroundDetails = ({
  id,
  patientName,
  dateOfRecording,
  summary,
  condensedSummary,
  timeOfRecording,
  durationOfRecording,
}: SearchResultItem) => {
  const durationInMinutes = Math.round(durationOfRecording / 60000);

  const getDuration = (duration: number) => {
    duration = Math.floor(duration / 60);
    const durationToString = `${String(duration)} minutes`;
    return duration > 1 ? durationToString : '< 1 min';
  };

  return (
    <div className="bg-gray-50 rounded-xl shadow-lg border border-gray-200 p-6 mx-6 my-6">
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6">
        <div>
          <h3 className="text-lg font-semibold text-gray-700">Last Visit</h3>
          <p className="mt-1 text-gray-600">
            {new Date(dateOfRecording).toLocaleDateString()}
          </p>
        </div>

        <div>
          <h3 className="text-lg font-semibold text-gray-700">Duration</h3>
          <p className="mt-1 text-gray-600">{getDuration(durationInMinutes)}</p>
        </div>
      </div>

      <div className="mb-6">
        <h2 className="text-2xl font-bold text-gray-800">Patient Background</h2>
        <p className="mt-2 text-gray-600 whitespace-pre-line">
          {condensedSummary}
        </p>
      </div>

      <div>
        <h3 className="text-lg font-semibold text-gray-700">
          Last Visit Summary
        </h3>
        <p className="mt-2 text-gray-600 whitespace-pre-line">{summary}</p>
      </div>
    </div>
  );
};

export default PatientBackgroundDetails;
