import { useState } from 'react';
import GoogleIcon from '../../assets/google-icon.webp';
import landingBg from '../../assets/landing-bg.png';
import logo from '../../assets/medscribe-logo.png';
import { useNavigate } from 'react-router-dom';
const SignUp = () => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (password !== confirmPassword) {
      alert('Passwords do not match');
      return;
    }
    console.log(name, email, password);
  };

  return (
    <div className="flex h-screen">
      {/* Left side - Sign Up Form */}
      <div className="w-1/3 px-[50px] flex flex-col">
        {/* Logo and Name */}
        <div className="mt-8 flex items-center gap-2">
          <img src={logo} alt="Medscribe Logo" className="w-16 h-16" />
          <span className="text-3xl font-semibold text-gray-800">Medscribe</span>
        </div>

        <div className="flex-1 flex items-center px-[50px]">
          <div className="flex flex-col gap-8 w-full">
            <h1 className="text-3xl font-semibold text-center text-gray-800">
              Create Your Account
            </h1>

            {/* Form */}
            <form onSubmit={handleSubmit} className="w-full">
              <div className="flex flex-col gap-4">
                <div className="relative">
                  <input
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    className="w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Full Name"
                    required
                  />
                </div>

                <div className="relative">
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Email"
                    required
                  />
                </div>

                <div className="relative">
                  <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Password"
                    required
                  />
                </div>

                <div className="relative">
                  <input
                    type="password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    className="w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Confirm Password"
                    required
                  />
                </div>

                <button
                  type="submit"
                  className="w-full py-3 mt-4 text-white text-base rounded-lg bg-blue-600 hover:bg-blue-700 transition-colors duration-200"
                >
                  Sign Up
                </button>
              </div>
            </form>

            {/* Divider */}
            <div className="flex items-center gap-4">
              <div className="flex-grow h-px bg-gray-300"></div>
              <span className="text-sm text-gray-500">OR</span>
              <div className="flex-grow h-px bg-gray-300"></div>
            </div>

            {/* Google Sign Up Button */}
            <button
              onClick={() => console.log('Sign up with Google')}
              className="w-full py-3 px-4 flex items-center justify-center gap-3 border border-gray-300 rounded-lg hover:bg-gray-50 hover:border-blue-500 transition-colors duration-200"
            >
              <img src={GoogleIcon} alt="Google" className="w-6 h-6" />
              <span className="text-gray-700">Sign up with Google</span>
            </button>

            {/* Sign in link */}
            <div className="text-center">
              <p className="text-gray-700">
                Already have an account?{' '}
                <button
                  onClick={() => navigate('/signin')}
                  className="text-blue-600 hover:text-blue-700 transition-colors duration-200"
                >
                  Sign in here
                </button>
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Right side - Background Image */}
      <div 
        className="w-2/3 h-full bg-cover bg-center bg-no-repeat"
        style={{ 
          backgroundImage: `url(${landingBg})`,
          backgroundSize: 'cover'
        }}
      />
    </div>
  );
};

export default SignUp;
