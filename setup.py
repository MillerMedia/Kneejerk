from setuptools import setup, find_packages

# Read version
with open('VERSION', 'r') as f:
    version = f.read().strip()

setup(
    name="kneejerk",
    version=version,
    description="Kneejerk - A tool for scanning environment variables from React websites",
    packages=find_packages(),
    install_requires=[
        'beautifulsoup4',
        'requests',
    ],
    entry_points={
        'console_scripts': [
            'kneejerk=kneejerk.scanner:main',
        ],
    },
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
)
