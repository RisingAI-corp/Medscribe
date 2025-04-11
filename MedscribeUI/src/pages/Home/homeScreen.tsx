import { Button } from '@mantine/core';
import PatientRegistryLayout from '../../components/PatientRegistry/patientRegistryLayout';
import PatientDashBoardLayout from '../../components/PaitentDashboard/patientDashBoardLayout';
import { Header } from '../../components/Header/header';
import { currentlySelectedPatientAtom } from '../../states/patientsAtom';
import PatientReception from '../../components/PatientReception/patientReceptionLayout';
import { useAtom } from 'jotai';

function HomeScreen() {
  const [patientId, setPatientId] = useAtom(currentlySelectedPatientAtom);
  return (
    <div className="flex flex-col flex-1 overflow-hidden h-screen">
      <div className="min-h-[30px]">
        <Header />
      </div>

      <div className="flex flex-1 flex-row overflow-hidden">
        <div className="flex flex-col w-[300px] border-r border-gray-300 h-full box-border overflow-y-auto">
          <div className="p-4 border-b border-gray-300 bg-white">
            <Button
              variant="outline"
              className="border-blue-500 text-blue-500 font-bold w-full"
              onClick={() => {
                setPatientId('');
              }}
            >
              + START A VISIT
            </Button>
          </div>
          <div className="flex-1 p-4">
            <PatientRegistryLayout />
          </div>
        </div>
        <div className="flex-1 h-full overflow-hidden relative">
          {/* Reception: always rendered, toggle visibility */}
          <div
            className={`${patientId ? 'hidden' : 'block'} w-full h-full overflow-y-auto`}
          >
            <PatientReception />
          </div>

          {/* Dashboard: always rendered, toggle visibility */}
          <div
            className={`${patientId ? 'block' : 'hidden'} w-full h-full overflow-y-auto`}
          >
            <PatientDashBoardLayout key={patientId || 'default'} />
          </div>
        </div>
      </div>
    </div>
  );
}

export default HomeScreen;
